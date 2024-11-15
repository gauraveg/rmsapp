package models

type LoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"gte=6,lte=15"`
}

type SessionToken struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

type LoginData struct {
	UserID       string `db:"id"`
	PasswordHash string `db:"password"`
	Role         string `db:"role"`
}

type UserCtx struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
	Role      string `json:"role"`
	Email     string `json:"email"`
}

type User struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Role       string        `json:"role"`
	Email      string        `json:"email"`
	Address    []AddressData `json:"address"`
	CreatedBy  string        `json:"createdBy"`
	CreatedAt  string        `json:"createdAt"`
	UpdatedBy  string        `json:"updatedBy"`
	UpdatedAt  string        `json:"updatedAt"`
	ArchivedAt *string       `json:"archivedAt"`
}

type SubAdminRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"gte=6,lte=15"`
}

type UserData struct {
	Name      string        `json:"name" validate:"required"`
	Email     string        `json:"email" validate:"email"`
	Password  string        `json:"password" validate:"gte=6,lte=15"`
	Role      string        `json:"role"`
	Addresses []AddressData `json:"addresses"`
}

type AddressData struct {
	Id         *string  `json:"id"`
	Address    string   `json:"address" validate:"required"`
	Latitude   *float64 `json:"latitude"`
	Longitude  *float64 `json:"longitude"`
	UserId     *string  `json:"userId"`
	CreatedAt  *string  `json:"createdAt"`
	ArchivedAt *string  `json:"archivedAt"`
}

type Restaurant struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	CreatedBy  string  `json:"createdBy"`
	CreatedAt  string  `json:"createdAt"`
	ArchivedAt *string `json:"archivedAt"`
}

type RestaurantsRequest struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type DishRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Dish struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Price        int     `json:"price"`
	RestaurantId string  `json:"restaurantId"`
	CreatedAt    string  `json:"createdAt"`
	ArchivedAt   *string `json:"archivedAt"`
}

type SessionData struct {
	Email      string  `json:"email"`
	ArchivedAt *string `json:"archivedAt"`
}
