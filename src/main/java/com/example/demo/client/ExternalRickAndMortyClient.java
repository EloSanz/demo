package com.example.demo.client;

import com.example.demo.dto.rickandmorty.RickAndMortyCharacterResponseDto;
import com.example.demo.dto.rickandmorty.RickAndMortyPageResponseDto;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.service.annotation.GetExchange;
import org.springframework.web.service.annotation.HttpExchange;

@HttpExchange("/character")
public interface ExternalRickAndMortyClient {

    /**
     * Get characters with optional filters and pagination. Generates: GET
     * /character/?page=1&name=rick&status=alive
     */
    @GetExchange
    RickAndMortyPageResponseDto<RickAndMortyCharacterResponseDto> getCharacters(
            @RequestParam(required = false) Integer page,
            @RequestParam(required = false) String name,
            @RequestParam(required = false) String status,
            @RequestParam(required = false) String species);

    @GetExchange("/{id}")
    RickAndMortyCharacterResponseDto getCharacterById(@PathVariable Long id);
}
