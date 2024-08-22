<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class GetTransactionsResponseCursor
{
    /**
     * $data
     *
     * @var array<WalletsTransaction> $data
     */
    #[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('array<\formance\stack\Models\Shared\WalletsTransaction>')]
    public array $data;

    /**
     *
     * @var ?bool $hasMore
     */
    #[\JMS\Serializer\Annotation\SerializedName('hasMore')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?bool $hasMore = null;

    /**
     *
     * @var ?string $next
     */
    #[\JMS\Serializer\Annotation\SerializedName('next')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $next = null;

    /**
     *
     * @var int $pageSize
     */
    #[\JMS\Serializer\Annotation\SerializedName('pageSize')]
    public int $pageSize;

    /**
     *
     * @var ?string $previous
     */
    #[\JMS\Serializer\Annotation\SerializedName('previous')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $previous = null;

    /**
     * @param  ?array<WalletsTransaction>  $data
     * @param  ?int  $pageSize
     * @param  ?bool  $hasMore
     * @param  ?string  $next
     * @param  ?string  $previous
     */
    public function __construct(?array $data = null, ?int $pageSize = null, ?bool $hasMore = null, ?string $next = null, ?string $previous = null)
    {
        $this->data = $data;
        $this->pageSize = $pageSize;
        $this->hasMore = $hasMore;
        $this->next = $next;
        $this->previous = $previous;
    }
}