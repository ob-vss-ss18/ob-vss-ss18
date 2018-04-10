package moviestore

import (
	"errors"
)

// A Moviestore interface.
type Moviestore interface {
	AddMovie(string, FSK) Serial
	AddUser(string, Age) UserID
	Rent(Serial, UserID) (User, Movie, error)
	Return(Serial) (User, Movie, error)
	RentedByUser(UserID) ([]Movie, error)
}

type moviestoreImpl struct {
	available  map[Serial]Movie
	rented     map[UserID][]Movie
	users      map[UserID]User
	nextSerial Serial
	nextUserID UserID
}

// NewMoviestore generates a Moviestore which contains available movies,
// users and a mapping
// between users and currently rented movies.
func NewMoviestore() Moviestore {
	ms := new(moviestoreImpl)
	ms.available = make(map[Serial]Movie)
	ms.users = make(map[UserID]User)
	ms.rented = make(map[UserID][]Movie)
	return ms
}

// AddMovie adds a movie to the available movies map which is part of
// the moviestoreImpl struct:
// 		available	map[Serial]Movie
// The serial will be generated and returned.
func (ms *moviestoreImpl) AddMovie(title string, fsk FSK) Serial {
	serial := ms.nextSerial
	ms.nextSerial++
	movie := Movie{title, fsk, serial}
	ms.available[serial] = movie
	return serial
}

// AddUser adds an user to the users map which is part of
// the moviestoreImpl struct:
// 		users	map[UserID]User
// The userid will be generated and returned.
func (ms *moviestoreImpl) AddUser(name string, age Age) UserID {
	userID := ms.nextUserID
	ms.nextUserID++
	user := User{name, age, userID}
	ms.users[userID] = user
	return userID
}

/*
Rent a movie.
If the user is in users and the movie is in available, the movie will
be removed from available and appended to the slice of rented movies
by this user. Therefore, the moviestoreImpl struct contains a field
	rented    map[UserID][]Movie
The following error cases are handled and will be returned as error
containing the following texts:

- user not found

- movie not available for rent

- user ist too young
*/
func (ms *moviestoreImpl) Rent(serial Serial, userID UserID) (User, Movie, error) {
	user, userPresent := ms.users[userID]
	if !userPresent {
		return user, Movie{}, errors.New("user not found")
	}
	movie, movieAvailable := ms.available[serial]
	if !movieAvailable {
		return user, movie, errors.New("movie not available for rent")
	}
	if !movie.AllowedAtAge(user.Age) {
		return user, movie, errors.New("user ist too young")
	}
	delete(ms.available, serial)
	usersRentedMovies := ms.rented[userID]
	ms.rented[userID] = append(usersRentedMovies, movie)
	return user, movie, nil
}

/*
Return a movie.
Searches in all slices of the rented map, does some "housekeeping,
and returns the user and movie if found.
The error case is "movie not found in rented movies".
 */
func (ms *moviestoreImpl) Return(serial Serial) (User, Movie, error) {
	for user, movies := range ms.rented {
		for i, movie := range movies {
			if movie.Serial == serial {
				movies[i] = movies[len(movies)-1] // Replace it with the last one.
				ms.rented[user] = movies[:len(movies)-1]
				ms.available[serial] = movie
				return ms.users[user], movie, nil
			}
		}
	}
	return User{}, Movie{}, errors.New("movie not found in rented movies")
}

/*
RentedByUser returns a slice of movies currently rented by the user.
Error case is "userID unknown"
Be aware that slices are returned by reference.
 */
func (ms *moviestoreImpl) RentedByUser(userID UserID) ([]Movie, error) {
	_, userPresent := ms.users[userID]
	if !userPresent {
		return nil, errors.New("userID unknown")
	}
	movies, moviesPresent := ms.rented[userID]
	if !moviesPresent {
		return nil, nil
	}
	copyOfMovies := append([]Movie(nil), movies...)
	return copyOfMovies, nil
}
