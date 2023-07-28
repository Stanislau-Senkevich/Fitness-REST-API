package handler

import "Fitness_REST_API/internal/entity"

type workoutsResponse struct {
	Workouts []*entity.Workout `json:"workouts"`
}

type usersResponse struct {
	Users []*entity.User `json:"users"`
}

type usersInfoResponse struct {
	UsersInfo []*entity.UserInfo
}

type partnershipsResponse struct {
	Partnerships []*entity.Partnership `json:"partnerships"`
}

type idResponse struct {
	Id int64 `json:"id"`
}
type requestIdResponse struct {
	RequestId int64 `json:"request_id"`
}
type workoutIdResponse struct {
	WorkoutId int64 `json:"workout_id"`
}
type partnershipIdResponse struct {
	PartnershipId int64 `json:"partnership_id"`
}
