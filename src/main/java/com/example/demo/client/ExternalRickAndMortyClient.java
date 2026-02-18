package com.example.demo.client;

import com.example.demo.dto.rickandmorty.RickAndMortyCharacterResponse;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.service.annotation.GetExchange;
import org.springframework.web.service.annotation.HttpExchange;

@HttpExchange("/character")
public interface ExternalRickAndMortyClient {

    @GetExchange("/{id}")
    RickAndMortyCharacterResponse getCharacterById(@PathVariable Long id);
}
