package cdl

type Categories struct {
	Rooms            *parentCategory
	UserIsChatMember *parentCategory
	User             *parentCategory
	FindMemberBy     *parentCategory
	ChatIDByMemberID *parentCategory
}

func (d *Dataloader) ConfigureDataloader() {
	d.categories = &Categories{
		Rooms:            d.newRoomsCategory(),
		UserIsChatMember: d.newUserIsChatMemberCategory(),
		User:             d.newUserCategory(),
		FindMemberBy:     d.newFindMemberByCategory(),
		ChatIDByMemberID: d.newChatIDByMemberIDCategory(),
	}
}

func (d *Dataloader) newChatIDByMemberIDCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.chatIDByMemberID
	return c
}

func (d *Dataloader) newFindMemberByCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.findMemberBy
	return c
}

func (d *Dataloader) newUserCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.user
	return c
}

func (d *Dataloader) newUserIsChatMemberCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.userIsChatMember
	return c
}

func (d *Dataloader) newRoomsCategory() *parentCategory {
	c := d.newParentCategory()
	c.LoadFn = c.rooms
	return c
}
