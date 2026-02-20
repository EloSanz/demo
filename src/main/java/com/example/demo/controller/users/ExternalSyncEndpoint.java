package com.example.demo.controller.users;

import com.example.demo.domain.User;
import com.example.demo.dto.users.UserResponseDto;
import com.example.demo.mapper.UserMapper;
import com.example.demo.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import java.util.List;
import java.util.stream.Collectors;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class ExternalSyncEndpoint extends BaseUserController {

        private final UserMapper userMapper;

        public ExternalSyncEndpoint(UserService userService, UserMapper userMapper) {
                super(userService);
                this.userMapper = userMapper;
        }

        @Operation(summary = "Fetch external users", description = "Retrieves users from the external JSONPlaceholder API.")
        @ApiResponse(responseCode = "200", description = "External users retrieved successfully")
        @GetMapping("/external")
        public ResponseEntity<List<UserResponseDto>> fetchExternalUsers() {
                List<UserResponseDto> response = userService.fetchUsersFromExternalApi().stream()
                                .map(userMapper::toResponse)
                                .collect(Collectors.toList());
                return ResponseEntity.ok(response);
        }

        @Operation(summary = "Sync external user", description = "Fetches a user from the external API and saves them to the local database.")
        @ApiResponses(value = {
                        @ApiResponse(responseCode = "201", description = "User synced and created locally"),
                        @ApiResponse(responseCode = "404", description = "User not found in external API")
        })
        @PostMapping("/sync/{externalUserId}")
        public ResponseEntity<UserResponseDto> syncUserFromExternal(
                        @Parameter(description = "ID of the user in the external system") @PathVariable Long externalUserId) {
                User syncedUser = userService.syncUserFromExternalApi(externalUserId);
                return ResponseEntity.status(HttpStatus.CREATED).body(userMapper.toResponse(syncedUser));
        }
}
