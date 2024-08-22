<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Models\Shared;
class V2GetLedgerInfoResponse
{
    /**
     * HTTP response content type for this operation
     *
     * @var string $contentType
     */
    public string $contentType;

    /**
     * HTTP response status code for this operation
     *
     * @var int $statusCode
     */
    public int $statusCode;

    /**
     * Raw HTTP response; suitable for custom response parsing
     *
     * @var \Psr\Http\Message\ResponseInterface $rawResponse
     */
    public \Psr\Http\Message\ResponseInterface $rawResponse;

    /**
     * OK
     *
     * @var ?Shared\V2LedgerInfoResponse $v2LedgerInfoResponse
     */
    public ?Shared\V2LedgerInfoResponse $v2LedgerInfoResponse = null;

    /**
     * @param  ?string  $contentType
     * @param  ?int  $statusCode
     * @param  ?\Psr\Http\Message\ResponseInterface  $rawResponse
     * @param  ?Shared\V2LedgerInfoResponse  $v2LedgerInfoResponse
     */
    public function __construct(?string $contentType = null, ?int $statusCode = null, ?\Psr\Http\Message\ResponseInterface $rawResponse = null, ?Shared\V2LedgerInfoResponse $v2LedgerInfoResponse = null)
    {
        $this->contentType = $contentType;
        $this->statusCode = $statusCode;
        $this->rawResponse = $rawResponse;
        $this->v2LedgerInfoResponse = $v2LedgerInfoResponse;
    }
}