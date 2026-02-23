package com.example.demo.dto.rickandmorty;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import java.util.List;
import lombok.Data;

@Data
@JsonIgnoreProperties(ignoreUnknown = true)
public class RickAndMortyPageResponseDto<T> {
    private InfoDto info;
    private List<T> results;

    @Data
    @JsonIgnoreProperties(ignoreUnknown = true)
    public static class InfoDto {
        private Integer count;
        private Integer pages;
        private String next;
        private String prev;
    }
}
