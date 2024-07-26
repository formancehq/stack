<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


enum ErrorsEnum: string
{
    case Internal = 'INTERNAL';
    case InsufficientFund = 'INSUFFICIENT_FUND';
    case Validation = 'VALIDATION';
    case Conflict = 'CONFLICT';
    case NoScript = 'NO_SCRIPT';
    case CompilationFailed = 'COMPILATION_FAILED';
    case MetadataOverride = 'METADATA_OVERRIDE';
    case NotFound = 'NOT_FOUND';
}
