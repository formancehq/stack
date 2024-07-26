<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Utils\SpeakeasyMetadata;
class V2CreateBulkRequest
{
    /**
     * $requestBody
     *
     * @var ?array<mixed> $requestBody
     */
    #[SpeakeasyMetadata('request:mediaType=application/json')]
    public ?array $requestBody = null;

    /**
     * Name of the ledger.
     *
     * @var string $ledger
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=ledger')]
    public string $ledger;

    public function __construct()
    {
        $this->requestBody = null;
        $this->ledger = '';
    }
}