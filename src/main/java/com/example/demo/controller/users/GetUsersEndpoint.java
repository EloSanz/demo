package com.example.demo.controller.users;

import com.example.demo.dto.users.UserResponseDto;
import com.example.demo.mapper.UserMapper;
import com.example.demo.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class GetUsersEndpoint extends BaseUserController {

    private final UserMapper userMapper;

    public GetUsersEndpoint(UserService userService, UserMapper userMapper) {
        super(userService);
        this.userMapper = userMapper;
    }

    @Operation(summary = "Get all users", description = "Retrieves a paginated list of all users.")
    @ApiResponse(responseCode = "200", description = "Page of users retrieved successfully")
    @GetMapping
    public ResponseEntity<Page<UserResponseDto>> getAllUsers(@ParameterObject Pageable pageable) {
        Page<UserResponseDto> response =
                userService.getAllUsers(pageable).map(userMapper::toResponse);
        return ResponseEntity.ok(response);
    }

    @Operation(
            summary = "Get user by ID",
            description = "Retrieves a specific user by their unique identifier.")
    @ApiResponses(
            value = {
                @ApiResponse(responseCode = "200", description = "User found"),
                @ApiResponse(responseCode = "404", description = "User not found")
            })
    @GetMapping("/{id}")
    public ResponseEntity<UserResponseDto> getUserById(
            @Parameter(description = "ID of the user to retrieve") @PathVariable Long id) {
        return ResponseEntity.ok(userMapper.toResponse(userService.getUserById(id)));
    }
}
