package com.example.demo.aspect;

import java.lang.annotation.ElementType;
import java.lang.annotation.Inherited;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * Annotation to enable automatic logging for a method or class. When applied to a class (like a
 * Controller), all its public methods will be logged. Logs method entry, arguments, return value,
 * and execution time. In case of exception, it logs the exception details.
 */
@Target({ElementType.TYPE, ElementType.METHOD})
@Retention(RetentionPolicy.RUNTIME)
@Inherited
public @interface Loggable {
    /** Whether to log the method arguments. */
    boolean logArgs() default true;

    /** Whether to log the return value. */
    boolean logResult() default true;
}
