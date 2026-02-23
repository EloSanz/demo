package com.example.demo.domain;

import lombok.Builder;
import lombok.Getter;

@Getter
@Builder
public class RickAndMortyPageInfo {
    private Integer count;
    private Integer pages;
    private String next;
    private String prev;
}
