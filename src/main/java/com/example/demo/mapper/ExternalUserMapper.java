package com.example.demo.mapper;

import com.example.demo.domain.User;
import com.example.demo.dto.users.UserResponseDto;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface ExternalUserMapper {

    User toDomain(UserResponseDto externalResponse);
}
