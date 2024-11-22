package models

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleSubAdmin Role = "sub-admin"
	RoleUser     Role = "user"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"gte=6,lte=15"`
}

type SessionToken struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

type LoginData struct {
	UserID       string `db:"id"`
	PasswordHash string `db:"password"`
	Role         Role   `db:"role"`
}

type UserCtx struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
	Role      Role   `json:"role"`
	Email     string `json:"email"`
}

type User struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Role      Role          `json:"role"`
	Email     string        `json:"email"`
	Address   []AddressData `json:"address"`
	CreatedBy string        `json:"createdBy"`
	CreatedAt string        `json:"createdAt"`
	UpdatedBy string        `json:"updatedBy"`
	UpdatedAt string        `json:"updatedAt"`
}

type SubAdminRequest struct {
	Name     string `json:"name" validate:"required,UserNameCheck"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"gte=6,lte=15"`
}

type UserData struct {
	Name      string        `json:"name" validate:"required,UserNameCheck"`
	Email     string        `json:"email" validate:"required,email"`
	Password  string        `json:"password" validate:"gte=6,lte=15"`
	Role      Role          `json:"role" validate:"oneof=admin sub-admin user"`
	Addresses []AddressData `json:"addresses" validate:"omitnil,dive"`
}

type UserSignUp struct {
	Name      string        `json:"name" validate:"required,UserNameCheck"`
	Email     string        `json:"email" validate:"required,email"`
	Password  string        `json:"password" validate:"gte=6,lte=15"`
	Addresses []AddressData `json:"addresses"`
}

type SignUpWithRole struct {
	UserSignUp
	Role Role
}

type AddressData struct {
	Id        *string  `json:"id"`
	Address   string   `json:"address" validate:"required,AddressCheck"`
	Latitude  *float64 `json:"latitude" validate:"omitnil,number"`
	Longitude *float64 `json:"longitude" validate:"omitnil,number"`
	UserId    *string  `json:"userId"`
	CreatedAt *string  `json:"createdAt"`
}

type Restaurant struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Address   string     `json:"address"`
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
	Dishes    []DishData `json:"dishes"`
	CreatedBy string     `json:"createdBy"`
	CreatedAt string     `json:"createdAt"`
}

type RestaurantsRequest struct {
	Name      string  `json:"name" validate:"required,UserNameCheck"`
	Address   string  `json:"address" validate:"required,AddressCheck"`
	Latitude  float64 `json:"latitude" validate:"omitnil,number"`
	Longitude float64 `json:"longitude" validate:"omitnil,number"`
}

type DishRequest struct {
	Name  string `json:"name" validate:"required,UserNameCheck"`
	Price int    `json:"price" validate:"required,number"`
}

type DishData struct {
	Id           string `json:"Id"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	RestaurantId string `json:"restaurantId"`
	CreatedAt    string `json:"createdAt"`
}

type Dish struct {
	Id             string `json:"Id"`
	Name           string `json:"name"`
	Price          int    `json:"price"`
	RestaurantId   string `json:"restaurantId"`
	RestaurantName string `json:"restaurantName"`
	CreatedAt      string `json:"createdAt"`
}

type SessionData struct {
	Email      string  `json:"email"`
	ArchivedAt *string `json:"archivedAt"`
}
