<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2ListTriggersResponseCursor
{
    /**
     * $data
     *
     * @var array<V2Trigger> $data
     */
    #[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('array<\formance\stack\Models\Shared\V2Trigger>')]
    public array $data;

    /**
     *
     * @var bool $hasMore
     */
    #[\JMS\Serializer\Annotation\SerializedName('hasMore')]
    public bool $hasMore;

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
     * @param  ?array<V2Trigger>  $data
     * @param  ?bool  $hasMore
     * @param  ?int  $pageSize
     * @param  ?string  $next
     * @param  ?string  $previous
     */
    public function __construct(?array $data = null, ?bool $hasMore = null, ?int $pageSize = null, ?string $next = null, ?string $previous = null)
    {
        $this->data = $data;
        $this->hasMore = $hasMore;
        $this->pageSize = $pageSize;
        $this->next = $next;
        $this->previous = $previous;
    }
}