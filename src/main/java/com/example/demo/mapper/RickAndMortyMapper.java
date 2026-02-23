package com.example.demo.mapper;

import com.example.demo.domain.RickAndMortyCharacter;
import com.example.demo.domain.RickAndMortyCharacterPage;
import com.example.demo.domain.RickAndMortyPageInfo;
import com.example.demo.dto.rickandmorty.RickAndMortyCharacterResponseDto;
import com.example.demo.dto.rickandmorty.RickAndMortyPageResponseDto;
import java.util.List;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface RickAndMortyMapper {
    RickAndMortyCharacter toDomain(RickAndMortyCharacterResponseDto response);

    RickAndMortyCharacterResponseDto toResponse(RickAndMortyCharacter domain);

    List<RickAndMortyCharacter> toDomainList(List<RickAndMortyCharacterResponseDto> response);

    List<RickAndMortyCharacterResponseDto> toResponseList(List<RickAndMortyCharacter> domain);

    RickAndMortyPageInfo toDomainInfo(RickAndMortyPageResponseDto.InfoDto info);

    RickAndMortyCharacterPage toDomainPage(
            RickAndMortyPageResponseDto<RickAndMortyCharacterResponseDto> response);

    RickAndMortyPageResponseDto.InfoDto fromDomainInfo(RickAndMortyPageInfo info);

    RickAndMortyPageResponseDto<RickAndMortyCharacterResponseDto> fromDomainPage(
            RickAndMortyCharacterPage domain);
}
