package com.example.demo.controller.rickAndMorty;

import static org.assertj.core.api.Assertions.assertThat;

import com.example.demo.BaseIntegrationTest;
import com.example.demo.dto.rickandmorty.RickAndMortyCharacterResponseDto;
import com.example.demo.util.RickAndMortyTestUtil;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.wiremock.spring.EnableWireMock;

@EnableWireMock
class RickAndMortyIntegrationTest extends BaseIntegrationTest {

    @Autowired private TestRestTemplate restTemplate;

    @Test
    void shouldReturnCharacter_WhenExternalApiReturnsSuccess() {
        // Given
        RickAndMortyTestUtil.stubRickAndMortyCharacter(
                1, RickAndMortyTestUtil.getRickSanchezJson());

        // When
        ResponseEntity<RickAndMortyCharacterResponseDto> response =
                restTemplate.getForEntity(
                        "/api/rickandmorty/characters/1", RickAndMortyCharacterResponseDto.class);

        // Then
        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);
        assertThat(response.getBody()).isNotNull();
        assertThat(response.getBody().getName()).isEqualTo("Rick Sanchez");
        assertThat(response.getBody().getStatus()).isEqualTo("Alive");
    }
}
