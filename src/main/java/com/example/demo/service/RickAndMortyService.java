package com.example.demo.service;

import com.example.demo.domain.RickAndMortyCharacter;
import com.example.demo.domain.RickAndMortyCharacterPage;

public interface RickAndMortyService {
    RickAndMortyCharacter getCharacterById(Long id);

    RickAndMortyCharacterPage getCharacters(
            Integer page, String name, String status, String species);
}
