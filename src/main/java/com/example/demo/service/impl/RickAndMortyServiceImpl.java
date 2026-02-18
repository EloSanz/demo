package com.example.demo.service.impl;

import com.example.demo.client.ExternalRickAndMortyClient;
import com.example.demo.dto.rickandmorty.RickAndMortyCharacterResponse;
import com.example.demo.service.RickAndMortyService;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class RickAndMortyServiceImpl implements RickAndMortyService {

    private final ExternalRickAndMortyClient client;

    @Override
    public RickAndMortyCharacterResponse getCharacterById(Long id) {
        return client.getCharacterById(id);
    }
}
