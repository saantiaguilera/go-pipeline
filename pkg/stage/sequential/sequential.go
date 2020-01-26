package sequential

// Run synchronously a number of workers. If one of them fails, the operation is aborted and the error returned.
func runSync(workers int, run func(index int) error) error {
	for i := 0; i < workers; i++ {
		err := run(i)

		if err != nil {
			return err
		}
	}
	return nil
}
