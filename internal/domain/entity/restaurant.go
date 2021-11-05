package entity

type Restaurant struct {
	Name string `json:"name"`
	MenuItems int `json:"menu_items"`
	Menus []Menu `json:"menu"`
}
