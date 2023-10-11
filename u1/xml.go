package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type Movie struct {
	Title       string `xml:"title,attr"`
	Genre       string `xml:"genre,attr"`
	ReleaseDate string `xml:"releaseDate,attr"`
}

type Movies struct {
	XMLName xml.Name `xml:"movies"`
	Movies  []Movie  `xml:"movie"`
}

func getMoviesByGenre(xmlFile string, genre string) ([]string, error) {
	var movieList []string

	xmlData, err := ioutil.ReadFile(xmlFile)
	if err != nil {
		return nil, err
	}

	var movies Movies
	err = xml.Unmarshal(xmlData, &movies)
	if err != nil {
		return nil, err
	}

	for _, movie := range movies.Movies {
		if movie.Genre == genre {
			movieList = append(movieList, movie.Title)
		}
	}

	return movieList, nil
}

func main() {

	xmlFile := flag.String("file", "herd.xml", "a string")
	genre := flag.String("genre", "crime", "a string")

	flag.Parse()

	//fmt.Println("file:", *xmlFile)
	//fmt.Println("genre:", *genre)

	movies, err := getMoviesByGenre(*xmlFile, *genre)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(movies) > 0 {
		fmt.Printf("Movies in the '%s' genre:\n", *genre)
		for _, movie := range movies {
			fmt.Println(movie)
		}
	} else {
		fmt.Printf("No movies found in the '%s' genre.\n", *genre)
	}
}
