package com.example.demo.service.impl;

import com.example.demo.client.ExternalRickAndMortyClient;
import com.example.demo.domain.RickAndMortyCharacter;
import com.example.demo.exception.rickandmorty.CharacterNotFoundException;
import com.example.demo.mapper.RickAndMortyMapper;
import com.example.demo.service.RickAndMortyService;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClientResponseException;

@Service
@RequiredArgsConstructor
public class RickAndMortyServiceImpl implements RickAndMortyService {

    private final ExternalRickAndMortyClient client;
    private final RickAndMortyMapper mapper;

    @Override
    public RickAndMortyCharacter getCharacterById(Long id) {
        try {
            return mapper.toDomain(client.getCharacterById(id));
        } catch (WebClientResponseException.NotFound ex) {
            throw new CharacterNotFoundException(id, ex);
        }
    }
}
