package cdl

func (c *parentCategory) getRequest(ptr chanPtr) *baseRequest {
	request, ok := c.Requests[ptr]
	if !ok { // если еще не создавали то надо паниковать
		panic("c.Requests not exists by" + ptr)
	}
	return request
}
