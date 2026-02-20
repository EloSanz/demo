package com.example.demo.controller.users;

import com.example.demo.domain.User;
import com.example.demo.dto.users.UserRequestDto;
import com.example.demo.dto.users.UserResponseDto;
import com.example.demo.mapper.UserMapper;
import com.example.demo.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import jakarta.validation.Valid;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
@Slf4j
public class CreateUserEndpoint extends BaseUserController {

    private final UserMapper userMapper;

    public CreateUserEndpoint(UserService userService, UserMapper userMapper) {
        super(userService);
        this.userMapper = userMapper;
    }

    @Operation(
            summary = "Create a new user",
            description = "Creates a new user in the local database. Email must be unique.")
    @ApiResponses(
            value = {
                @ApiResponse(responseCode = "201", description = "User created successfully"),
                @ApiResponse(
                        responseCode = "400",
                        description = "Invalid input or email already exists")
            })
    @PostMapping
    public ResponseEntity<UserResponseDto> createUser(@Valid @RequestBody UserRequestDto request) {
        User domainUser = userMapper.toDomain(request);
        User createdUser = userService.createUser(domainUser);
        return ResponseEntity.status(HttpStatus.CREATED).body(userMapper.toResponse(createdUser));
    }
}
