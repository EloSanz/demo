package com.example.demo.controller.users;

import com.example.demo.domain.User;
import com.example.demo.dto.users.UserRequestDto;
import com.example.demo.dto.users.UserResponseDto;
import com.example.demo.mapper.UserMapper;
import com.example.demo.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class UpdateUserEndpoint extends BaseUserController {

    private final UserMapper userMapper;

    public UpdateUserEndpoint(UserService userService, UserMapper userMapper) {
        super(userService);
        this.userMapper = userMapper;
    }

    @Operation(summary = "Update user", description = "Updates an existing user's information.")
    @ApiResponses(
            value = {
                @ApiResponse(responseCode = "200", description = "User updated successfully"),
                @ApiResponse(responseCode = "400", description = "Invalid input"),
                @ApiResponse(responseCode = "404", description = "User not found")
            })
    @PutMapping("/{id}")
    public ResponseEntity<UserResponseDto> updateUser(
            @Parameter(description = "ID of the user to update") @PathVariable Long id,
            @Valid @RequestBody UserRequestDto request) {
        User domainUser = userMapper.toDomain(request);
        User updatedUser = userService.updateUser(id, domainUser);
        return ResponseEntity.ok(userMapper.toResponse(updatedUser));
    }
}
