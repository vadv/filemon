package server

import (
	"github.com/hpcloud/tail"
	"sync"
)

type fileProcessor struct {
	lock   *sync.Mutex
	tailer *tail.Tail
	// map[key]LinkToObject
	commands map[string]command
}

func newFileProcessor(filename string) (*fileProcessor, error) {

	t, err := tail.TailFile(filename, tail.Config{
		MaxLineSize: 8 * 1024, // MaxLineSize
		Follow:      true,     // tail -f
		ReOpen:      true,     // Reopen recreated files (tail -F)
		MustExist:   true,     // Fail early if the file does not exist
		Poll:        true,     // Poll for file changes instead of using inotify
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
		//Logger: &log.Logger{},
	})
	if err != nil {
		return nil, err
	}

	result := &fileProcessor{
		lock:     &sync.Mutex{},
		tailer:   t,
		commands: make(map[string]command, 0),
	}

	go result.run()

	return result, nil
}

func (f *fileProcessor) run() {
	for line := range f.tailer.Lines {

		f.lock.Lock()
		commands := f.commands
		f.lock.Unlock()

		for _, cmd := range commands {
			cmd.Process(line.Text)
		}
	}
}

func (f *fileProcessor) registerCommand(key, name, expr string, args []string) error {
	cmd, err := newCommand(name, expr, args)
	if err != nil {
		return err
	}
	f.commands[key] = cmd
	return nil
}
