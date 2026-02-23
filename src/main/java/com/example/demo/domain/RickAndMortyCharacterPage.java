package com.example.demo.domain;

import java.util.List;
import lombok.Builder;
import lombok.Getter;

@Getter
@Builder
public class RickAndMortyCharacterPage {
    private RickAndMortyPageInfo info;
    private List<RickAndMortyCharacter> results;
}
