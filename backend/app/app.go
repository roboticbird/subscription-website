package app

import (
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"strconv"
	"subscription-website/backend/tasks"
	"time"
)

type App struct {
	Router   *mux.Router
	Store    *sessions.CookieStore
	Database *sql.DB
}

type UserSession struct {
	UserDetails   *User
	Authenticated bool
}

func (app *App) Setup() {
	// set up session store
	app.Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
	gob.Register(UserSession{})

	// set up router
	app.Router.
		Methods("GET").
		Path("/forbidden").
		HandlerFunc(app.forbiddenHandler)
	app.Router.
		Methods("GET").
		Path("/userInfo").
		HandlerFunc(app.userInfoHandler)
	app.Router.
		Methods("POST").
		Path("/login").
		HandlerFunc(app.loginHandler)
	app.Router.
		Methods("POST").
		Path("/logout").
		HandlerFunc(app.logoutHandler)
	app.Router.
		Methods("GET").
		Path("/products").
		HandlerFunc(app.productsHandler)
	app.Router.
		Methods("POST").
		Path("/product").
		HandlerFunc(app.productHandler)
	app.Router.
		Methods("POST").
		Path("/newSubscription").
		HandlerFunc(app.newSubscriptionHandler)
	app.Router.
		Methods("POST").
		Path("/updateSubscription").
		HandlerFunc(app.updateSubscriptionHandler)
	app.Router.
		Methods("GET").
		Path("/subscriptions").
		HandlerFunc(app.subscriptionsHandler)
	// for testing only
	app.Router.
		Methods("POST").
		Path("/cronUpdate").
		HandlerFunc(app.cronUpdateHandler)
}

func (app *App) forbiddenHandler(w http.ResponseWriter, r *http.Request) {
	session, err := app.Store.Get(r, "session.id")
	if err != nil {
		log.Printf("Session store error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Forbidden\n"))
}

func (app *App) userInfoHandler(w http.ResponseWriter, r *http.Request) {
	session, err := app.Store.Get(r, "session.id")
	if err != nil {
		log.Printf("Session store error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := app.getCurrentUser(session)

	if !app.checkAuthentication(user) {
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user.UserDetails); err != nil {
		log.Printf("Json Encoding Error error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	type LoginData struct {
		Email    string
		Password string
	}

	session, _ := app.Store.Get(r, "session.id")
	decoder := json.NewDecoder(r.Body)
	var login LoginData
	err := decoder.Decode(&login)

	if err != nil {
		log.Printf("Login error: %s\n", err)
		http.Error(w, "Please pass the data as JSON", http.StatusBadRequest)
		return
	}

	user, err := app.getUser(login.Email)
	if err != nil {
		log.Printf("Login error for %s: %s", login.Email, err.Error())
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	if login.Password != user.Password {
		log.Printf("Login error for %s: %s", login.Email, err.Error())
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	sessionUser := &UserSession{
		UserDetails:   user,
		Authenticated: true,
	}
	session.Values["user"] = sessionUser
	session.Options = &sessions.Options{SameSite: http.SameSiteStrictMode}
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Unable to save session %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Logged in successfully\n"))
}

func (app *App) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := app.Store.Get(r, "session.id")
	session.Values["user"] = UserSession{}
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		log.Printf("Unable to save session %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Logout in successfully\n"))
}

func (app *App) getProduct(product_id string) (*Product, error) {

	product := Product{}
	err := app.Database.QueryRow("SELECT * FROM products WHERE id = ?", product_id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Duration)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (app *App) productHandler(w http.ResponseWriter, r *http.Request) {
	type ProductData struct {
		ProductId int
	}
	decoder := json.NewDecoder(r.Body)
	var prodData ProductData
	err := decoder.Decode(&prodData)
	if err != nil {
		log.Printf("Product error: %s\n", err)
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}

	product, err := app.getProduct(strconv.Itoa(prodData.ProductId))
	if err != nil {
		log.Printf("Session store error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Printf("Json Encoding Error error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *App) productsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := app.Database.Query("SELECT * FROM products")
	if err != nil {
		log.Printf("Error reading database: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	products := []Product{}
	for rows.Next() {
		product := Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Duration)
		if err != nil {
			log.Printf("Error reading database: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		log.Printf("Error reading database: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("Error encoding data: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *App) newSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	type NewSubData struct {
		ProductId int
	}

	session, err := app.Store.Get(r, "session.id")
	if err != nil {
		log.Printf("Session store error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := app.getCurrentUser(session)

	if !app.checkAuthentication(user) {
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var newSub NewSubData
	err = decoder.Decode(&newSub)
	if err != nil {
		log.Printf("New subscription error: %s\n", err)
		http.Error(w, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}

	subscription, err := app.getSubscription(user.UserDetails, strconv.Itoa(newSub.ProductId))
	if subscription != nil {
		log.Printf("User trying to buy duplicate subscription: user=%d, product=%d, status=%s\n", user.UserDetails.ID, newSub.ProductId, subscription.Status)
		w.WriteHeader(http.StatusNotModified)
		return
	}
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error reading database: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = app.Database.Exec(fmt.Sprintf("INSERT INTO subscriptions (user_id, product_id, start_date, status, nextStatus) VALUES(%d, %d, NOW(), 'ACTIVE', 'QUEUED')", user.UserDetails.ID, newSub.ProductId))
	if err != nil {
		log.Printf("Error writing to database: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *App) subscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := app.Store.Get(r, "session.id")
	if err != nil {
		log.Printf("Session store error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := app.getCurrentUser(session)

	if !app.checkAuthentication(user) {
		err = session.Save(r, w)
		if err != nil {
			log.Printf("Session store error: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	subscriptions, err := app.getSubscriptions(user.UserDetails)
	if err != nil {
		log.Printf("Error reading database: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(subscriptions); err != nil {
		log.Printf("Error encoding data: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *App) updateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	type UpdateSubData struct {
		ProductId int
		Action    string
	}

	session, err := app.Store.Get(r, "session.id")
	if err != nil {
		log.Printf("Session store error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := app.getCurrentUser(session)

	if !app.checkAuthentication(user) {
		err = session.Save(r, w)
		if err != nil {
			log.Printf("Session store error: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var updateSub UpdateSubData
	err = decoder.Decode(&updateSub)
	if err != nil {
		log.Printf("Update subscription error: %s\n", err)
		http.Error(w, "Please pass the data as JSON", http.StatusBadRequest)
		return
	}

	subscription, err := app.getSubscription(user.UserDetails, strconv.Itoa(updateSub.ProductId))
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error reading database: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if subscription == nil || subscription.Status == "CANCELLED" {
		log.Printf("Trying to update invalid subscription\n")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	query := ""
	switch updateSub.Action {
	case "QUEUE":
		query = fmt.Sprintf("UPDATE subscriptions SET nextStatus = 'QUEUED' WHERE id = %d", subscription.ID)
	case "PAUSE":
		query = fmt.Sprintf("UPDATE subscriptions SET nextStatus = 'PAUSED' WHERE id = %d", subscription.ID)
	case "CANCEL":
		query = fmt.Sprintf("UPDATE subscriptions SET status = 'CANCELLED' WHERE id = %d", subscription.ID)
	default:
		log.Printf("Trying to update to invalid status code: %s\n", updateSub.Action)
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	_, err = app.Database.Exec(query)
	if err != nil {
		log.Printf("Error writing to database: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

// for testing only
func (app *App) cronUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type CronUpdateData struct {
		PrivateToken string
	}
	type Sub struct {
		ID         int
		ProductID  int
		UserID     int
		StartDate  time.Time
		Status     string
		NextStatus string
	}
	token := "top-secret"
	decoder := json.NewDecoder(r.Body)
	var cronUpdate CronUpdateData
	err := decoder.Decode(&cronUpdate)
	if err != nil {
		log.Printf("Cron Update Error: %s\n", err)
		http.Error(w, "Please pass the data as JSON", http.StatusBadRequest)
		return
	}
	if cronUpdate.PrivateToken != token {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	tasks.OrderUpdate(app.Database)
	w.WriteHeader(http.StatusOK)
	return

}

func (app *App) getCurrentUser(s *sessions.Session) UserSession {
	val := s.Values["user"]
	user := UserSession{}
	user, ok := val.(UserSession)
	if !ok {
		return UserSession{Authenticated: false}
	}
	return user
}

func (app *App) getUser(email string) (*User, error) {
	user := User{}
	err := app.Database.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (app *App) checkAuthentication(user UserSession) bool {
	if auth := user.Authenticated; !auth {
		log.Printf("User trying to access secret without auth\n")
		return false
	}
	return true
}

func (app *App) getSubscription(user *User, product_id string) (*Subscription, error) {
	subscription := Subscription{}
	err := app.Database.QueryRow("SELECT sub.id, sub.start_date, sub.status, sub.nextStatus, "+
		"prd.id, prd.name, prd.description, prd.price, prd.duration "+
		"FROM subscriptions as sub "+
		"LEFT JOIN products as prd on prd.id = sub.product_id "+
		"WHERE sub.product_id="+product_id+" "+
		"AND sub.user_id="+strconv.Itoa(user.ID)+" "+
		"AND sub.status != 'CANCELLED' "+
		"AND sub.status != 'EXPIRED'").Scan(&subscription.ID,
		&subscription.StartDate, &subscription.Status, &subscription.NextStatus,
		&subscription.Product.ID, &subscription.Product.Name,
		&subscription.Product.Description, &subscription.Product.Price,
		&subscription.Product.Duration)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Error reading database: %s", err)
		return nil, err
	}
	return &subscription, err
}

func (app *App) getSubscriptions(user *User) ([]Subscription, error) {
	rows, err := app.Database.Query("SELECT sub.id, sub.start_date, sub.status, sub.nextStatus, " +
		"prd.id, prd.name, prd.description, prd.price, prd.duration " +
		"FROM subscriptions as sub " +
		"LEFT JOIN products as prd on prd.id = sub.product_id " +
		"WHERE sub.user_id=" + strconv.Itoa(user.ID) + " " +
		"AND sub.status != 'CANCELLED' " +
		"AND sub.status != 'EXPIRED'")

	if err != nil {
		log.Printf("Error reading database: %s", err)
		return nil, err
	}
	defer rows.Close()
	subscriptions := []Subscription{}
	for rows.Next() {
		subscription := Subscription{}
		err := rows.Scan(&subscription.ID,
			&subscription.StartDate, &subscription.Status, &subscription.NextStatus,
			&subscription.Product.ID, &subscription.Product.Name,
			&subscription.Product.Description, &subscription.Product.Price,
			&subscription.Product.Duration)
		if err != nil {
			log.Printf("Error reading database: %s", err)
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	err = rows.Err()
	return subscriptions, err
}
