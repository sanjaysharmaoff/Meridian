package main

import "math/rand"

// Predefined routes for drivers (used for the gRPC Streaming module)
// (these are San Francisco routes, get these coordinates from Google Maps for example and build a custom route if you want)
var PredefinedRoutes = [][][]float64{
	{
		{18.9432, 72.8311},
		{18.9448, 72.8325},
		{18.9464, 72.8338},
		{18.9481, 72.8345},
	},
	{
		{18.9347, 72.8258},
		{18.9361, 72.8264},
		{18.9372, 72.8269},
		{18.9364, 72.8275},
		{18.9380, 72.8282},
		{18.9395, 72.8289},
		{18.9408, 72.8294},
		{18.9412, 72.8311},
		{18.9401, 72.8318},
		{18.9390, 72.8312},
		{18.9384, 72.8330},
	},
	{
		{18.9439, 72.8247},
		{18.9422, 72.8253},
		{18.9410, 72.8258},
		{18.9404, 72.8271},
		{18.9397, 72.8286},
		{18.9389, 72.8299},
		{18.9378, 72.8295},
		{18.9366, 72.8290},
	},
	{
		{18.9366, 72.8290},
		{18.9378, 72.8295},
		{18.9389, 72.8299},
		{18.9397, 72.8286},
		{18.9404, 72.8271},
		{18.9422, 72.8253},
		{18.9439, 72.8247},
		{18.9410, 72.8258},
	},
}

func GenerateRandomPlate() string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	plate := ""
	for i := 0; i < 3; i++ {
		plate += string(letters[rand.Intn(len(letters))])
	}

	return plate
}
