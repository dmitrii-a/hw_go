package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	result := in
	for _, stage := range stages {
		if stage == nil {
			continue
		}
		result = executeStage(stage, result, done)
	}
	return result
}

func executeStage(stage Stage, in In, done In) Out {
	result := make(Bi)
	go func() {
		defer close(result)
		in := stage(in)
		for {
			select {
			case <-done:
				return
			case item, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case result <- item:
				}
			}
		}
	}()
	return result
}
