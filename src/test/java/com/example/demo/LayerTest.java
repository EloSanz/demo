package com.example.demo;

import static com.tngtech.archunit.library.Architectures.layeredArchitecture;

import com.tngtech.archunit.core.domain.JavaClasses;
import com.tngtech.archunit.core.importer.ClassFileImporter;
import com.tngtech.archunit.core.importer.ImportOption;
import com.tngtech.archunit.lang.ArchRule;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

public class LayerTest {

    private static final String CONTROLLER = "controller";
    private static final String SERVICE = "service";
    private static final String REPOSITORY = "repository";
    private static final String DOMAIN = "domain";
    private static final String CLIENT = "client";
    private static final String DTO = "dto";
    private static final String CONFIG = "config";

    private JavaClasses javaClasses;

    @BeforeEach
    public void init() {
        this.javaClasses =
                new ClassFileImporter()
                        .withImportOption(ImportOption.Predefined.DO_NOT_INCLUDE_JARS)
                        .withImportOption(ImportOption.Predefined.DO_NOT_INCLUDE_TESTS)
                        .importPackages("com.example.demo");
    }

    @Test
    void validatingAccessBetweenLayersOfCode() {
        ArchRule mainLayer =
                layeredArchitecture()
                        .consideringOnlyDependenciesInLayers()
                        .layer(CONTROLLER)
                        .definedBy("..controller..")
                        .layer(SERVICE)
                        .definedBy("..service..")
                        .layer(REPOSITORY)
                        .definedBy("..repository..")
                        .layer(DOMAIN)
                        .definedBy("..domain..")
                        .layer(CLIENT)
                        .definedBy("..client..")
                        .layer(DTO)
                        .definedBy("..dto..")
                        .layer(CONFIG)
                        .definedBy("..config..")
                        .whereLayer(CONTROLLER)
                        .mayNotBeAccessedByAnyLayer()
                        .whereLayer(SERVICE)
                        .mayOnlyBeAccessedByLayers(CONTROLLER, CONFIG)
                        .whereLayer(REPOSITORY)
                        .mayOnlyBeAccessedByLayers(SERVICE, CONFIG)
                        .whereLayer(CLIENT)
                        .mayOnlyBeAccessedByLayers(SERVICE, CONFIG)
                        .whereLayer(DOMAIN)
                        .mayOnlyBeAccessedByLayers(
                                REPOSITORY, SERVICE, CLIENT, CONTROLLER, CONFIG, DTO) // Domain
                        // is
                        // core
                        .whereLayer(DTO)
                        .mayOnlyBeAccessedByLayers(
                                CONTROLLER, SERVICE, CLIENT, CONFIG) // DTOs are data
                        // transfer objects
                        .withOptionalLayers(true);

        mainLayer.check(this.javaClasses);
    }
}
