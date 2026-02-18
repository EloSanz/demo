package com.example.demo.repository;

import com.example.demo.domain.User;
import java.util.Optional;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

/**
 * Repository interface for User entity. Database-agnostic using Spring Data JPA.
 *
 * <p>You can easily swap H2 for PostgreSQL, MySQL, etc. by changing: - application.yml datasource
 * config - build.gradle database driver dependency
 */
@Repository
public interface UserRepository extends JpaRepository<User, Long> {

    Optional<User> findByEmail(String email);

    boolean existsByEmail(String email);
}
