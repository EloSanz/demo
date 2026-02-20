package com.example.demo.mapper;

import com.example.demo.domain.User;
import com.example.demo.entity.UserEntity;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.MappingTarget;

@Mapper(componentModel = "spring")
public interface UserEntityMapper {

    UserEntity toEntity(User domain);

    User toDomain(UserEntity entity);

    @Mapping(target = "id", ignore = true)
    void updateEntityFromDomain(User domain, @MappingTarget UserEntity entity);
}
