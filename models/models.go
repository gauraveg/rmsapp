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
	UserID       string `db:"userid"`
	PasswordHash string `db:"password"`
	Role         string `db:"role"`
}

type UserCtx struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
	Role      string `json:"role"`
}

type User struct {
	UserID     string        `json:"userId"`
	Name       string        `json:"name"`
	Role       string        `json:"role"`
	Email      string        `json:"email"`
	Address    []AddressData `json:"address"`
	CreatedBy  string        `json:"createdBy"`
	CreatedAt  string        `json:"createdAt"`
	Updatedby  string        `json:"updatedBy"`
	UpdatedAt  string        `json:"updatedAt"`
	ArchivedAt *string       `json:"archivedAt"`
}

type SubAdminRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"gte=6,lte=15"`
}

type UserData struct {
	Name     string        `json:"name" validate:"required"`
	Email    string        `json:"email" validate:"email"`
	Password string        `json:"password" validate:"gte=6,lte=15"`
	Role     string        `json:"role"`
	Address  []AddressData `json:"address"`
}

type AddressData struct {
	AddressId   *string  `json:"addressId"`
	AddressLine string   `json:"addressline" validate:"required"`
	Latitude    *float64 `json:"latitude" validate:"required"`
	Longitude   *float64 `json:"longitude" validate:"required"`
	User_Id     *string  `json:"user_id"`
	CreatedAt   *string  `json:"createdat"`
	ArchivedAt  *string  `json:"archivedat"`
}

type Restaurant struct {
	RestaurantId string  `json:"restaurantId"`
	Name         string  `json:"name"`
	AddressLine  string  `json:"addressline"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	CreatedBy    string  `json:"createdBy"`
	CreatedAt    string  `json:"createdAt"`
	ArchivedAt   *string `json:"arcivedAt"`
}

type RestaurantsRequest struct {
	Name        string  `json:"name"`
	AddressLine string  `json:"addressline"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
