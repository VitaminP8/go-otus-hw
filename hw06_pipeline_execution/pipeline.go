package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	// начиная с входного канала, проходимся по всем стейджам и связываем их
	out := in
	for _, stage := range stages {
		out = doStage(out, done, stage)
	}

	return out
}

func doStage(in In, done In, stage Stage) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for value := range in {
			select {
			case <-done:
				return
			default:
				// обрабатываем данные с помощью stage (передаем канал, содержащий только value и учитывающий done)
				stageOut := stage(wrapValueToChan(value))

				select {
				case <-done:
					return
				default:
					out <- <-stageOut
				}
			}
		}
	}()
	return out
}

func wrapValueToChan(value any) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		out <- value
	}()

	return out
}
