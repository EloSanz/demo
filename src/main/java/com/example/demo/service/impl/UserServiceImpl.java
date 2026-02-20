package com.example.demo.service.impl;

import com.example.demo.client.ExternalUserClient;
import com.example.demo.domain.User;
import com.example.demo.dto.users.UserResponseDto;
import com.example.demo.entity.UserEntity;
import com.example.demo.exception.users.UserAlreadyExistsException;
import com.example.demo.exception.users.UserNotFoundException;
import com.example.demo.mapper.ExternalUserMapper;
import com.example.demo.mapper.UserEntityMapper;
import com.example.demo.repository.UserRepository;
import com.example.demo.service.UserService;
import java.util.List;
import java.util.stream.Collectors;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

/** Implementation of UserService. Contains business logic for User operations. */
@Service
@RequiredArgsConstructor
@Slf4j
@Transactional
public class UserServiceImpl implements UserService {

    private final UserRepository userRepository;
    private final ExternalUserClient externalUserClient;
    private final UserEntityMapper entityMapper;
    private final ExternalUserMapper externalMapper;

    @Override
    @Transactional(readOnly = true)
    public Page<User> getAllUsers(Pageable pageable) {
        return userRepository.findAll(pageable).map(entityMapper::toDomain);
    }

    @Override
    @Transactional(readOnly = true)
    public User getUserById(Long id) {
        UserEntity entity =
                userRepository.findById(id).orElseThrow(() -> new UserNotFoundException(id));
        return entityMapper.toDomain(entity);
    }

    @Override
    public User createUser(User domainObject) {
        if (userRepository.existsByEmail(domainObject.getEmail())) {
            throw new UserAlreadyExistsException(domainObject.getEmail());
        }

        UserEntity entity = entityMapper.toEntity(domainObject);
        UserEntity savedEntity = userRepository.save(entity);

        return entityMapper.toDomain(savedEntity);
    }

    @Override
    public User updateUser(Long id, User domainObject) {
        UserEntity existingEntity =
                userRepository.findById(id).orElseThrow(() -> new UserNotFoundException(id));

        entityMapper.updateEntityFromDomain(domainObject, existingEntity);

        UserEntity updatedEntity = userRepository.save(existingEntity);

        return entityMapper.toDomain(updatedEntity);
    }

    @Override
    public void deleteUser(Long id) {
        if (!userRepository.existsById(id)) {
            throw new UserNotFoundException(id);
        }
        userRepository.deleteById(id);
    }

    @Override
    @Transactional(readOnly = true)
    public List<User> fetchUsersFromExternalApi() {
        return externalUserClient.getAllUsers().stream()
                .map(externalMapper::toDomain)
                .collect(Collectors.toList());
    }

    @Override
    @Transactional
    public User syncUserFromExternalApi(Long externalUserId) {
        UserResponseDto externalResponse = externalUserClient.getUserById(externalUserId);
        User incomingUser = externalMapper.toDomain(externalResponse);

        UserEntity existingEntity =
                userRepository
                        .findByEmail(incomingUser.getEmail())
                        .orElseGet(() -> entityMapper.toEntity(incomingUser));

        entityMapper.updateEntityFromDomain(incomingUser, existingEntity);
        UserEntity savedEntity = userRepository.save(existingEntity);

        return entityMapper.toDomain(savedEntity);
    }
}
