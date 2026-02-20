package com.example.demo.mapper;

import com.example.demo.domain.User;
import com.example.demo.dto.users.UserRequestDto;
import com.example.demo.dto.users.UserResponseDto;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface UserMapper {

    User toDomain(UserRequestDto request);

    UserResponseDto toResponse(User domain);
}
