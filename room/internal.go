package room

func (r *Room) promoteNewOwner() {
	for _, user := range r.users {
		user.Promote()
		r.owner = user
		break
	}
}
