package cdl

type Categories struct {
	Rooms            *ParentCategory
	UserIsChatMember *ParentCategory
	User             *ParentCategory
}

func (d *Dataloader) ConfigureDataloader() {
	d.Categories = &Categories{
		Rooms:            d.NewRoomsCategory(),
		UserIsChatMember: d.NewUserIsChatMemberCategory(),
		User:             d.NewUserCategory(),
	}
}

func (d *Dataloader) NewRoomsCategory() *ParentCategory {
	c := d.NewParentCategory()
	c.LoadFn = c.rooms
	return c
}

func (d *Dataloader) NewUserIsChatMemberCategory() *ParentCategory {
	c := d.NewParentCategory()
	c.LoadFn = c.userIsChatMember
	return c
}

func (d *Dataloader) NewUserCategory() *ParentCategory {
	c := d.NewParentCategory()
	c.LoadFn = c.user
	return c
}
