package com.example.demo.controller.rickAndMorty;

import static com.github.tomakehurst.wiremock.client.WireMock.aResponse;
import static com.github.tomakehurst.wiremock.client.WireMock.get;
import static com.github.tomakehurst.wiremock.client.WireMock.stubFor;
import static com.github.tomakehurst.wiremock.client.WireMock.urlEqualTo;
import static org.assertj.core.api.Assertions.assertThat;

import com.example.demo.AbstractIntegrationTest;
import com.example.demo.dto.rickandmorty.RickAndMortyCharacterResponse;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.test.context.TestPropertySource;
import org.wiremock.spring.ConfigureWireMock;
import org.wiremock.spring.EnableWireMock;

@EnableWireMock(@ConfigureWireMock(port = 0))
@TestPropertySource(properties = {
                "external.api.rickandmorty.base-url=http://localhost:${wiremock.server.port}/api"
})
class RickAndMortyIntegrationTest extends AbstractIntegrationTest {

        @Autowired
        private TestRestTemplate restTemplate;

        @Test
        void shouldReturnCharacter_WhenExternalApiReturnsSuccess() {
                // Given
                String externalResponse = """
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

                stubFor(get(urlEqualTo("/api/character/1"))
                                .willReturn(aResponse()
                                                .withStatus(200)
                                                .withHeader("Content-Type", MediaType.APPLICATION_JSON_VALUE)
                                                .withBody(externalResponse)));

                // When
                ResponseEntity<RickAndMortyCharacterResponse> response = restTemplate
                                .getForEntity("/api/rickandmorty/characters/1", RickAndMortyCharacterResponse.class);

                // Then
                assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);
                assertThat(response.getBody()).isNotNull();
                assertThat(response.getBody().getName()).isEqualTo("Rick Sanchez");
                assertThat(response.getBody().getStatus()).isEqualTo("Alive");
        }
}
