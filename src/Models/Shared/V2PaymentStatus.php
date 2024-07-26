<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


enum V2PaymentStatus: string
{
    case Pending = 'PENDING';
    case Active = 'ACTIVE';
    case Terminated = 'TERMINATED';
    case Failed = 'FAILED';
    case Succeeded = 'SUCCEEDED';
    case Cancelled = 'CANCELLED';
}
