package com.example.demo.util;

import static com.github.tomakehurst.wiremock.client.WireMock.aResponse;
import static com.github.tomakehurst.wiremock.client.WireMock.get;
import static com.github.tomakehurst.wiremock.client.WireMock.stubFor;
import static com.github.tomakehurst.wiremock.client.WireMock.urlEqualTo;

import org.springframework.http.MediaType;

public class RickAndMortyTestUtil {

    private RickAndMortyTestUtil() {
        // Private constructor for utility class
    }

    public static void stubRickAndMortyCharacter(int id, String responseBody) {
        stubFor(get(urlEqualTo("/api/character/" + id))
                .willReturn(aResponse()
                        .withStatus(200)
                        .withHeader("Content-Type", MediaType.APPLICATION_JSON_VALUE)
                        .withBody(responseBody)));
    }

    public static String getRickSanchezJson() {
        return """
                {
                    "id": 1,
                    "name": "Rick Sanchez",
                    "status": "Alive",
                    "species": "Human",
                    "type": "",
                    "gender": "Male",
                    "image": "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
                    "url": "https://rickandmortyapi.com/api/character/1",
                    "created": "2017-11-04T18:48:46.250Z"
                }
                """;
    }
}
