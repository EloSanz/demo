package com.example.demo.mapper;

import com.example.demo.domain.User;
import com.example.demo.dto.users.UserResponseDto;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface ExternalUserMapper {

    @Mapping(target = "id", ignore = true)
    User toDomain(UserResponseDto externalResponse);
}
