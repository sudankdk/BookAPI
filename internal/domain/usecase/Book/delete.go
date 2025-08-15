package book

func (c *CreateBookUsecase) Delete(id string) error {
	return c.repo.Delete(id)
}
