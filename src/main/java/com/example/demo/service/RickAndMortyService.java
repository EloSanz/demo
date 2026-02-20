package com.example.demo.service;

import com.example.demo.domain.RickAndMortyCharacter;

public interface RickAndMortyService {
    RickAndMortyCharacter getCharacterById(Long id);
}
