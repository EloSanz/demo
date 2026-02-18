package com.example.demo.service;

import com.example.demo.dto.rickandmorty.RickAndMortyCharacterResponse;

public interface RickAndMortyService {
    RickAndMortyCharacterResponse getCharacterById(Long id);
}
