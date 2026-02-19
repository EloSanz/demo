package com.example.demo.architecture;

import static com.tngtech.archunit.lang.syntax.ArchRuleDefinition.classes;

import com.tngtech.archunit.core.importer.ImportOption;
import com.tngtech.archunit.junit.AnalyzeClasses;
import com.tngtech.archunit.junit.ArchTest;
import com.tngtech.archunit.lang.ArchRule;

@AnalyzeClasses(packages = "com.example.demo", importOptions = ImportOption.DoNotIncludeTests.class)
public class ExceptionArchitectureTest {

    @ArchTest
    static final ArchRule exceptions_should_reside_in_exception_package =
            classes()
                    .that()
                    .haveSimpleNameEndingWith("Exception")
                    .should()
                    .resideInAPackage("..exception..");

    @ArchTest
    static final ArchRule exceptions_should_be_exceptions =
            classes()
                    .that()
                    .resideInAPackage("..exception..")
                    .should()
                    .beAssignableTo(Exception.class)
                    .orShould()
                    .beAssignableTo(RuntimeException.class);

    @ArchTest
    static final ArchRule services_should_throw_domain_exceptions =
            classes()
                    .that()
                    .resideInAPackage("..service..")
                    .and()
                    .areNotInterfaces()
                    .should()
                    .dependOnClassesThat()
                    .resideInAPackage("..exception..");
}
