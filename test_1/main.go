package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type ProcessResult struct {
	FileName    string
	RowCount    int
	ProcessTime time.Duration
	Error       error
}

type FileJob struct {
	FilePath string
	FileNum  int
}

type ProgressTracker struct {
	mu        sync.Mutex
	total     int
	completed int
	failed    int
}

func (pt *ProgressTracker) Update(fileName string, success bool) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	
	pt.completed++
	if !success {
		pt.failed++
	}
	
	status := "✓"
	if !success {
		status = "✗"
	}
	
	fmt.Printf("[%d/%d] %s %s\n", pt.completed, pt.total, status, fileName)
}

type ConcurrentProcessor struct {
	workerCount int
	results     []ProcessResult
	resultsMu   sync.Mutex
	tracker     *ProgressTracker
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewProcessor(workerCount int) *ConcurrentProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	return &ConcurrentProcessor{
		workerCount: workerCount,
		results:     make([]ProcessResult, 0),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (cp *ConcurrentProcessor) ProcessFiles(filePaths []string) []ProcessResult {
	cp.tracker = &ProgressTracker{total: len(filePaths)}

	jobs := make(chan FileJob, len(filePaths))
	results := make(chan ProcessResult, len(filePaths))

	var wg sync.WaitGroup
	for i := 0; i < cp.workerCount; i++ {
		wg.Add(1)
		go cp.worker(jobs, results, &wg)
	}

	// Send jobs to workers
	go func() {
		defer close(jobs)
		for i, path := range filePaths {
			select {
			case jobs <- FileJob{FilePath: path, FileNum: i + 1}:
			case <-cp.ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		select {
		case <-cp.ctx.Done():
			// If context was cancelled, stop processing results
			return cp.results
		default:
			cp.resultsMu.Lock()
			cp.results = append(cp.results, result)
			cp.resultsMu.Unlock()

			cp.tracker.Update(result.FileName, result.Error == nil)
		}
	}

	return cp.results
}

func (cp *ConcurrentProcessor) worker(jobs <-chan FileJob, results chan<- ProcessResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		select {
		case <-cp.ctx.Done():
			// Context cancelled, exit worker
			return
		default:
			// Process the job normally
			result := cp.processFile(job.FilePath)

			// Send result to channel
			select {
			case results <- result:
			case <-cp.ctx.Done():
				// Context cancelled, exit worker
				return
			}
		}
	}
}

func (cp *ConcurrentProcessor) processFile(filePath string) ProcessResult {
	start := time.Now()
	result := ProcessResult{FileName: filepath.Base(filePath)}

	file, err := os.Open(filePath)
	if err != nil {
		result.Error = fmt.Errorf("open failed: %w", err)
		return result
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	rowCount := 0
	for {
		// Check for context cancellation periodically
		select {
		case <-cp.ctx.Done():
			result.Error = fmt.Errorf("processing cancelled: %w", cp.ctx.Err())
			return result
		default:
			// Continue processing
		}

		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Error = fmt.Errorf("read error at row %d: %w", rowCount+1, err)
			return result
		}

		if len(record) == 0 {
			result.Error = fmt.Errorf("empty record at row %d", rowCount+1)
			return result
		}

		// Simulate processing
		time.Sleep(1 * time.Millisecond)
		rowCount++
	}

	result.RowCount = rowCount
	result.ProcessTime = time.Since(start)
	return result
}

// Cancel stops all processing
func (cp *ConcurrentProcessor) Cancel() {
	cp.cancel()
}

// WithContext sets a custom context for the processor
func (cp *ConcurrentProcessor) WithContext(ctx context.Context) *ConcurrentProcessor {
	cp.cancel() // Cancel the current context
	cp.ctx, cp.cancel = context.WithCancel(ctx)
	return cp
}

func (cp *ConcurrentProcessor) PrintSummary() {
	fmt.Println("\n" + "==========================================================")
	fmt.Println("Processing Summary")
	fmt.Println("==========================================================")
	
	totalRows, totalTime, success := 0, time.Duration(0), 0
	
	for _, r := range cp.results {
		if r.Error == nil {
			fmt.Printf("✓ %s: %d rows in %v\n", r.FileName, r.RowCount, r.ProcessTime)
			totalRows += r.RowCount
			totalTime += r.ProcessTime
			success++
		} else {
			fmt.Printf("✗ %s: %v\n", r.FileName, r.Error)
		}
	}
	
	fmt.Println("==========================================================")
	fmt.Printf("Files: %d | Success: %d | Failed: %d\n", len(cp.results), success, len(cp.results)-success)
	fmt.Printf("Total Rows: %d | Avg Time: %v\n", totalRows, totalTime/time.Duration(len(cp.results)))
	fmt.Println("==========================================================")
}

func CreateSampleFiles(dir string, count int) ([]string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	
	paths := make([]string, count)
	
	for i := 0; i < count; i++ {
		path := filepath.Join(dir, fmt.Sprintf("sample_%d.csv", i+1))
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		
		w := csv.NewWriter(file)
		w.Write([]string{"ID", "Name", "Email", "Age", "City"})
		
		rows := 100 + (i * 50)
		for j := 0; j < rows; j++ {
			w.Write([]string{
				fmt.Sprintf("%d", j+1),
				fmt.Sprintf("User_%d_%d", i+1, j+1),
				fmt.Sprintf("user%d_%d@example.com", i+1, j+1),
				fmt.Sprintf("%d", 20+(j%50)),
				fmt.Sprintf("City_%d", (j%10)+1),
			})
		}
		
		w.Flush()
		file.Close()
		paths[i] = path
	}
	
	return paths, nil
}

func CalculateOptimalWorkers(fileCount int) int {
	// Strategi 1: Berdasarkan CPU cores
	cpuCount := runtime.NumCPU()
	
	// Strategi 2: Berdasarkan jumlah file
	// Tidak perlu worker lebih banyak dari file
	if fileCount < cpuCount {
		return fileCount
	}
	
	// Strategi 3: Max workers = 2x CPU cores (good for I/O bound tasks)
	maxWorkers := cpuCount * 2
	
	// Pilih yang paling efisien
	if fileCount < maxWorkers {
		return fileCount
	}
	
	return maxWorkers
}

func main() {
	fmt.Println("Concurrent CSV File Processor")
	fmt.Println("==========================================================\n")

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initial worker count - will be dynamically adjusted
	workerCount := 4

	// OPSI 1: Generate sample files otomatis
	useGenerated := false

	var files []string
	var err error

	if useGenerated {
		fileCount, dir := 10, "./csv_files"
		fmt.Printf("Creating %d sample files...\n", fileCount)
		files, err = CreateSampleFiles(dir, fileCount)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Created files in %s\n\n", dir)
		defer os.RemoveAll(dir) // Auto cleanup
	} else {
		// OPSI 2: Gunakan CSV files Anda sendiri
		// Ganti path sesuai lokasi CSV files Anda
		files = []string{
			"./data/file1.csv",
			"./data/file2.csv",
			"./data/file3.csv",
			// Tambah file lainnya...
		}
		fmt.Printf("Processing %d existing files...\n\n", len(files))
	}

	// Use dynamic worker count
	workerCount = CalculateOptimalWorkers(len(files))
	fmt.Printf("Processing with %d workers...\n\n", workerCount)
	processor := NewProcessor(workerCount).WithContext(ctx)

	// Goroutine to handle cancellation signals
	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived signal %v, cancelling processing...\n", sig)
		processor.Cancel()
	}()

	start := time.Now()
	processor.ProcessFiles(files)

	// Clean up
	cancel()
	signal.Stop(sigChan)
	close(sigChan)

	processor.PrintSummary()
	fmt.Printf("Total Time: %v\n\n", time.Since(start))

	fmt.Println("Done!")
}