package com.example.demo.controller.rickAndMorty;

import com.example.demo.dto.rickandmorty.RickAndMortyCharacterResponseDto;
import com.example.demo.dto.rickandmorty.RickAndMortyPageResponseDto;
import com.example.demo.mapper.RickAndMortyMapper;
import com.example.demo.service.RickAndMortyService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/rickandmorty/characters")
@RequiredArgsConstructor
@Tag(name = "Rick and Morty", description = "Operations related to Rick and Morty characters")
public class GetCharacterEndpoint {

    private final RickAndMortyService service;
    private final RickAndMortyMapper mapper;

    @Operation(
            summary = "Search characters",
            description = "Get characters with filtering and pagination.")
    @GetMapping
    public ResponseEntity<RickAndMortyPageResponseDto<RickAndMortyCharacterResponseDto>>
            searchCharacters(
                    @RequestParam(required = false) Integer page,
                    @RequestParam(required = false) String name,
                    @RequestParam(required = false) String status,
                    @RequestParam(required = false) String species) {
        return ResponseEntity.ok(
                mapper.fromDomainPage(service.getCharacters(page, name, status, species)));
    }

    @Operation(
            summary = "Get character by ID",
            description = "Fetches a character from the external Rick and Morty API.")
    @GetMapping("/{id}")
    public ResponseEntity<RickAndMortyCharacterResponseDto> getCharacter(@PathVariable Long id) {
        return ResponseEntity.ok(mapper.toResponse(service.getCharacterById(id)));
    }
}
