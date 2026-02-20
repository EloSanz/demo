package com.example.demo.mapper;

import com.example.demo.domain.RickAndMortyCharacter;
import com.example.demo.dto.rickandmorty.RickAndMortyCharacterResponseDto;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface RickAndMortyMapper {
    RickAndMortyCharacter toDomain(RickAndMortyCharacterResponseDto response);

    RickAndMortyCharacterResponseDto toResponse(RickAndMortyCharacter domain);
}
