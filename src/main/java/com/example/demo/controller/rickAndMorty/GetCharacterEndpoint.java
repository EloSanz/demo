package com.example.demo.controller.rickAndMorty;

import com.example.demo.dto.rickandmorty.RickAndMortyCharacterResponse;
import com.example.demo.service.RickAndMortyService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/rickandmorty/characters")
@RequiredArgsConstructor
@Tag(name = "Rick and Morty", description = "Operations related to Rick and Morty characters")
public class GetCharacterEndpoint {

    private final RickAndMortyService service;

    @Operation(
            summary = "Get character by ID",
            description = "Fetches a character from the external Rick and Morty API.")
    @GetMapping("/{id}")
    public ResponseEntity<RickAndMortyCharacterResponse> getCharacter(@PathVariable Long id) {
        return ResponseEntity.ok(service.getCharacterById(id));
    }
}
